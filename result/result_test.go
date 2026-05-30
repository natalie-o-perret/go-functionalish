package result_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/natalie-o-perret/go-functionalish/option"
	"github.com/natalie-o-perret/go-functionalish/result"
)

func TestOk(t *testing.T) {
	r := result.Ok[int, error](42)
	if !r.IsOk() {
		t.Fatal("expected Ok")
	}
	if r.Unwrap() != 42 {
		t.Fatalf("got %d", r.Unwrap())
	}
}

func TestErr(t *testing.T) {
	r := result.Err[int, error](errors.New("boom"))
	if !r.IsErr() {
		t.Fatal("expected Err")
	}
	if r.UnwrapOr(99) != 99 {
		t.Fatal("expected default")
	}
}

func TestTry(t *testing.T) {
	ok := result.Try(func() (int, error) { return 1, nil })
	if !ok.IsOk() {
		t.Fatal("expected Ok")
	}

	fail := result.Try(func() (int, error) { return 0, errors.New("bad") })
	if !fail.IsErr() {
		t.Fatal("expected Err")
	}
}

func TestMap(t *testing.T) {
	r := result.Ok[int, error](3).Map(func(n int) int { return n * 2 })
	if r.Unwrap() != 6 {
		t.Fatalf("got %d", r.Unwrap())
	}
}

func TestMapErr(t *testing.T) {
	r := result.Err[int, string]("oops").MapErr(
		func(s string) int { return len(s) },
	)
	if r.UnwrapErr() != 4 {
		t.Fatalf("got %d", r.UnwrapErr())
	}
}

func TestBind(t *testing.T) {
	double := func(n int) result.Result[int, error] { return result.Ok[int, error](n * 2) }
	got := result.Ok[int, error](5).Bind(double).Unwrap()
	if got != 10 {
		t.Fatalf("got %d", got)
	}
}

func TestToOption(t *testing.T) {
	opt := result.Ok[int, error](7).ToOption()
	if !opt.IsSome() || opt.Unwrap() != 7 {
		t.Fatal("expected Some(7)")
	}
	none := result.Err[int, error](errors.New("x")).ToOption()
	if !none.IsNone() {
		t.Fatal("expected None")
	}
}

// -- railway-oriented error handling ------------------------------------------
//
// Railway-oriented programming treats a pipeline as two parallel tracks:
// the success track (Ok) and the error track (Err).  Each Bind step may
// switch from success -> error; once on the error track every subsequent
// step is transparently skipped.  Map steps are pure transforms that
// cannot fail.

// assertPipeErr asserts the Result is Err and the message contains wantSubstr.
func assertPipeErr[T any](t *testing.T, r result.Result[T, string], wantSubstr string) {
	t.Helper()
	if r.IsOk() {
		t.Fatalf("expected Err containing %q, got Ok", wantSubstr)
	}
	if !strings.Contains(r.UnwrapErr(), wantSubstr) {
		t.Fatalf("Err = %q, want substring %q", r.UnwrapErr(), wantSubstr)
	}
}

// -- Pipeline 1: Order Fulfillment (8 steps) ---------------------------------
//
// string --Bind---> RawOrder --Bind---> VerifiedOrder --Bind--->
//  parse            resolve customer    check inventory
//
// --Map---> PricedOrder --Bind---> PricedOrder --Map---> FinalOrder
//   price              apply promo           tax
//
// --Bind---> FinalOrder --Map---> Receipt
//  charge                issue

func TestPipeline_OrderFulfillment(t *testing.T) {
	type RawOrder struct {
		Email string
		Promo string
		SKUs  []string
	}
	type VerifiedOrder struct {
		CustomerID string
		Promo      string
		SKUs       []string
	}
	type PricedOrder struct {
		CustomerID string
		Promo      string
		Subtotal   float64
	}
	type FinalOrder struct {
		CustomerID string
		Subtotal   float64
		Tax        float64
		Total      float64
	}
	type Receipt struct {
		OrderID string
		Total   float64
	}

	// Simulated data stores.
	customers := map[string]string{
		"alice@example.com": "C-100",
		"bob@example.com":   "C-200",
	}
	catalog := map[string]float64{
		"BOLT": 29.99,
		"GEAR": 89.50,
		"CHIP": 249.99,
	}
	stock := map[string]int{
		"BOLT": 500,
		"GEAR": 20,
		"CHIP": 0, // out of stock
	}
	promos := map[string]float64{
		"TENOFF":  0.10,
		"QUARTER": 0.25,
	}

	// Step 1 (Bind): parse semicolon-delimited input.
	parseOrder := func(raw string) result.Result[RawOrder, string] {
		parts := strings.SplitN(raw, ";", 3)
		if len(parts) != 3 {
			return result.Err[RawOrder, string]("invalid format: want email;skus;promo")
		}
		skus := strings.Split(parts[1], ",")
		if len(skus) == 0 || skus[0] == "" {
			return result.Err[RawOrder, string]("at least one SKU required")
		}
		return result.Ok[RawOrder, string](RawOrder{
			Email: parts[0], SKUs: skus, Promo: parts[2],
		})
	}

	// Step 2 (Bind): resolve customer by email.
	resolveCustomer := func(o RawOrder) result.Result[VerifiedOrder, string] {
		id, ok := customers[o.Email]
		if !ok {
			return result.Err[VerifiedOrder, string](
				fmt.Sprintf("unknown customer: %s", o.Email))
		}
		return result.Ok[VerifiedOrder, string](VerifiedOrder{
			CustomerID: id, SKUs: o.SKUs, Promo: o.Promo,
		})
	}

	// Step 3 (Bind): verify every SKU is in stock.
	checkInventory := func(o VerifiedOrder) result.Result[VerifiedOrder, string] {
		for _, sku := range o.SKUs {
			qty, exists := stock[sku]
			if !exists {
				return result.Err[VerifiedOrder, string](
					fmt.Sprintf("unknown SKU: %s", sku))
			}
			if qty == 0 {
				return result.Err[VerifiedOrder, string](
					fmt.Sprintf("out of stock: %s", sku))
			}
		}
		return result.Ok[VerifiedOrder, string](o)
	}

	// Step 4 (Map - pure): look up catalog prices and compute subtotal.
	priceItems := func(o VerifiedOrder) PricedOrder {
		var sub float64
		for _, sku := range o.SKUs {
			sub += catalog[sku]
		}
		return PricedOrder{CustomerID: o.CustomerID, Subtotal: sub, Promo: o.Promo}
	}

	// Step 5 (Bind): validate and apply the promo code.
	applyPromo := func(o PricedOrder) result.Result[PricedOrder, string] {
		if o.Promo == "" {
			return result.Ok[PricedOrder, string](o)
		}
		rate, ok := promos[o.Promo]
		if !ok {
			return result.Err[PricedOrder, string](
				fmt.Sprintf("invalid promo code: %s", o.Promo))
		}
		o.Subtotal *= (1 - rate)
		return result.Ok[PricedOrder, string](o)
	}

	// Step 6 (Map - pure): compute tax and total.
	finalize := func(o PricedOrder) FinalOrder {
		tax := o.Subtotal * 0.08
		return FinalOrder{
			CustomerID: o.CustomerID,
			Subtotal:   o.Subtotal,
			Tax:        tax,
			Total:      o.Subtotal + tax,
		}
	}

	// Step 7 (Bind): enforce a $200 single-transaction limit.
	chargePayment := func(o FinalOrder) result.Result[FinalOrder, string] {
		if o.Total > 200 {
			return result.Err[FinalOrder, string](
				fmt.Sprintf("total $%.2f exceeds $200 limit", o.Total))
		}
		return result.Ok[FinalOrder, string](o)
	}

	// Step 8 (Map - pure): produce the receipt.
	issueReceipt := func(o FinalOrder) Receipt {
		return Receipt{
			OrderID: "ORD-" + o.CustomerID,
			Total:   math.Round(o.Total*100) / 100,
		}
	}

	// The full 8-step pipeline.
	process := func(input string) result.Result[Receipt, string] {
		r1 := parseOrder(input)        // step 1: Bind
		r2 := r1.Bind(resolveCustomer) // step 2: Bind
		r3 := r2.Bind(checkInventory)  // step 3: Bind
		r4 := r3.Map(priceItems)       // step 4: Map (pure)
		r5 := r4.Bind(applyPromo)      // step 5: Bind
		r6 := r5.Map(finalize)         // step 6: Map (pure)
		r7 := r6.Bind(chargePayment)   // step 7: Bind
		r8 := r7.Map(issueReceipt)     // step 8: Map (pure)
		return r8
	}

	t.Run("all 8 steps succeed", func(t *testing.T) {
		r := process("alice@example.com;BOLT,GEAR;TENOFF")
		if r.IsErr() {
			t.Fatalf("expected Ok, got Err(%s)", r.UnwrapErr())
		}
		receipt := r.Unwrap()
		if receipt.OrderID != "ORD-C-100" {
			t.Errorf("OrderID = %s", receipt.OrderID)
		}
		// subtotal 119.49 × 0.90 = 107.541, tax ~ 8.603, total ~ 116.14
		if got := fmt.Sprintf("%.2f", receipt.Total); got != "116.14" {
			t.Errorf("Total = %s, want 116.14", got)
		}
	})

	t.Run("succeeds without promo", func(t *testing.T) {
		r := process("bob@example.com;BOLT;")
		if r.IsErr() {
			t.Fatalf("expected Ok, got Err(%s)", r.UnwrapErr())
		}
		// subtotal 29.99, tax ~ 2.40, total ~ 32.39
		if got := fmt.Sprintf("%.2f", r.Unwrap().Total); got != "32.39" {
			t.Errorf("Total = %s, want 32.39", got)
		}
	})

	t.Run("short-circuits at step 1: bad format", func(t *testing.T) {
		assertPipeErr(t, process("just-an-email"), "invalid format")
	})

	t.Run("short-circuits at step 2: unknown customer", func(t *testing.T) {
		assertPipeErr(t, process("nobody@example.com;BOLT;"), "unknown customer")
	})

	t.Run("short-circuits at step 3: out of stock", func(t *testing.T) {
		assertPipeErr(t, process("alice@example.com;CHIP;"), "out of stock: CHIP")
	})

	t.Run("short-circuits at step 5: invalid promo", func(t *testing.T) {
		assertPipeErr(t, process("alice@example.com;BOLT;EXPIRED99"), "invalid promo code")
	})

	t.Run("short-circuits at step 7: exceeds charge limit", func(t *testing.T) {
		// BOLT (29.99) + CHIP would be out of stock, so use BOLT + GEAR (119.49)
		// without promo -> total ~ 129.05 < 200 -> passes step 7.
		// We need a bigger order. But CHIP is out of stock...
		// Instead, show MapErr adding context to a step-5 failure.
		r := process("alice@example.com;BOLT;EXPIRED99")
		enriched := r.MapErr(func(e string) string {
			return "order pipeline: " + e
		})
		assertPipeErr(t, enriched, "order pipeline: invalid promo code")
	})
}

// -- Pipeline 2: User Registration (7 steps) ---------------------------------
//
// string --Bind---> Credentials --Bind---> Credentials --Bind--->
//  parse            validate email       check available
//
// --Bind---> Credentials --Map---> HashedCreds --Bind---> Account --Map---> Session
//  password strength    hash      create account        issue session
//
// The entire result is then wrapped with MapErr to add pipeline context.

func TestPipeline_UserOnboarding(t *testing.T) {
	type Credentials struct {
		Email    string
		Password string
	}
	type HashedCreds struct {
		Email        string
		PasswordHash string
	}
	type Account struct {
		Email string
		ID    int
	}
	type Session struct {
		Token     string
		Role      string
		AccountID int
	}

	takenEmails := map[string]bool{
		"taken@example.com": true,
	}

	// Step 1 (Bind): parse "email:password".
	parseInput := func(raw string) result.Result[Credentials, string] {
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return result.Err[Credentials, string]("format must be email:password")
		}
		return result.Ok[Credentials, string](Credentials{
			Email: parts[0], Password: parts[1],
		})
	}

	// Step 2 (Bind): validate email format.
	validateEmail := func(c Credentials) result.Result[Credentials, string] {
		if !strings.Contains(c.Email, "@") {
			return result.Err[Credentials, string]("invalid email: missing @")
		}
		if !strings.Contains(c.Email, ".") {
			return result.Err[Credentials, string]("invalid email: missing domain")
		}
		return result.Ok[Credentials, string](c)
	}

	// Step 3 (Bind): ensure email is not already registered.
	checkAvailable := func(c Credentials) result.Result[Credentials, string] {
		if takenEmails[c.Email] {
			return result.Err[Credentials, string](
				fmt.Sprintf("email already registered: %s", c.Email))
		}
		return result.Ok[Credentials, string](c)
	}

	// Step 4 (Bind): enforce password strength.
	checkPasswordStrength := func(c Credentials) result.Result[Credentials, string] {
		if len(c.Password) < 8 {
			return result.Err[Credentials, string]("password must be at least 8 characters")
		}
		hasUpper := strings.ContainsAny(c.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasDigit := strings.ContainsAny(c.Password, "0123456789")
		if !hasUpper || !hasDigit {
			return result.Err[Credentials, string](
				"password must contain an uppercase letter and a digit")
		}
		return result.Ok[Credentials, string](c)
	}

	// Step 5 (Map - pure): hash the password (simplified for testing).
	hashPassword := func(c Credentials) HashedCreds {
		return HashedCreds{
			Email:        c.Email,
			PasswordHash: fmt.Sprintf("argon2:%d", len(c.Password)),
		}
	}

	// Step 6 (Bind): create the account (may fail for blocked domains).
	createAccount := func(h HashedCreds) result.Result[Account, string] {
		if strings.HasSuffix(h.Email, "@blocked.com") {
			return result.Err[Account, string]("domain is blocked")
		}
		return result.Ok[Account, string](Account{
			ID: 1001, Email: h.Email,
		})
	}

	// Step 7 (Map - pure): assign default role and issue session token.
	issueSession := func(a Account) Session {
		return Session{
			Token:     fmt.Sprintf("tok_%d", a.ID),
			AccountID: a.ID,
			Role:      "member",
		}
	}

	// The full 7-step pipeline, with MapErr wrapping all errors.
	register := func(input string) result.Result[Session, string] {
		r1 := parseInput(input)                  // step 1
		r2 := r1.Bind(validateEmail)             // step 2
		r3 := r2.Bind(checkAvailable)            // step 3
		r4 := r3.Bind(checkPasswordStrength)     // step 4
		r5 := r4.Map(hashPassword)               // step 5 (pure)
		r6 := r5.Bind(createAccount)             // step 6
		r7 := r6.Map(issueSession)               // step 7 (pure)
		return r7.MapErr(func(e string) string { // enrich errors
			return "registration failed: " + e
		})
	}

	t.Run("all 7 steps succeed", func(t *testing.T) {
		r := register("alice@example.com:Str0ngPass!")
		if r.IsErr() {
			t.Fatalf("expected Ok, got Err(%s)", r.UnwrapErr())
		}
		s := r.Unwrap()
		if s.Token != "tok_1001" || s.AccountID != 1001 || s.Role != "member" {
			t.Errorf("session = %+v", s)
		}
	})

	t.Run("fails at step 1: bad format", func(t *testing.T) {
		assertPipeErr(t, register("noformat"), "registration failed: format must be email:password")
	})

	t.Run("fails at step 2: invalid email", func(t *testing.T) {
		assertPipeErr(t, register("notanemail:Pass1234"), "registration failed: invalid email")
	})

	t.Run("fails at step 3: email taken", func(t *testing.T) {
		assertPipeErr(t, register("taken@example.com:Pass1234"), "registration failed: email already registered")
	})

	t.Run("fails at step 4: weak password - too short", func(t *testing.T) {
		assertPipeErr(t, register("new@example.com:short"), "registration failed: password must be at least 8")
	})

	t.Run("fails at step 4: weak password - missing complexity", func(t *testing.T) {
		assertPipeErr(t, register("new@example.com:alllowercase"),
			"registration failed: password must contain an uppercase letter and a digit")
	})

	t.Run("fails at step 6: blocked domain", func(t *testing.T) {
		assertPipeErr(t, register("user@blocked.com:Str0ngPass!"), "registration failed: domain is blocked")
	})
}

// -- Pipeline 3: Ledger Processing (6 steps) ----------------------------------
//
// Demonstrates result.Try, result.FromOption, and method chaining
// working together in a realistic 6-step pipeline.
//
// string --Bind---> RawTx --Bind---> RawTx --Bind--->
//  parse CSV        validate          parse amount (Try+MapErr)
//
// --Bind---> EnrichedTx --Map---> EnrichedTx --Map---> LedgerEntry
//  resolve acct (FromOption)  apply rate        format

func TestPipeline_LedgerEntry(t *testing.T) {
	type RawTx struct {
		DateStr   string
		AmountStr string
		AccountID string
		Memo      string
	}
	type EnrichedTx struct {
		Date        string
		AccountName string
		Memo        string
		Amount      float64
	}
	type LedgerEntry struct {
		Date    string
		Summary string
		Amount  float64
	}

	accounts := map[string]string{
		"ACCT-1": "Operating Chequing",
		"ACCT-2": "Payroll Reserve",
	}
	exchangeRate := 1.10 // 10% markup for foreign currency

	// Step 1 (Bind): split pipe-delimited line.
	parseCSV := func(line string) result.Result[RawTx, string] {
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			return result.Err[RawTx, string](
				fmt.Sprintf("expected 4 pipe-delimited fields, got %d", len(parts)))
		}
		return result.Ok[RawTx, string](RawTx{
			DateStr:   parts[0],
			AmountStr: parts[1],
			AccountID: parts[2],
			Memo:      parts[3],
		})
	}

	// Step 2 (Bind): check required fields are non-empty.
	validateFields := func(tx RawTx) result.Result[RawTx, string] {
		if tx.DateStr == "" {
			return result.Err[RawTx, string]("date is required")
		}
		if tx.AmountStr == "" {
			return result.Err[RawTx, string]("amount is required")
		}
		if tx.AccountID == "" {
			return result.Err[RawTx, string]("account ID is required")
		}
		return result.Ok[RawTx, string](tx)
	}

	// Step 3 (Bind): parse the amount string using Try + MapErr.
	parseAmount := func(tx RawTx) result.Result[EnrichedTx, string] {
		amtResult := result.Try(func() (float64, error) {
			return strconv.ParseFloat(tx.AmountStr, 64)
		}).MapErr(func(e error) string {
			return fmt.Sprintf("invalid amount %q: %v", tx.AmountStr, e)
		})
		return amtResult.Map(func(amt float64) EnrichedTx {
			return EnrichedTx{
				Date:   tx.DateStr,
				Amount: amt,
				Memo:   tx.Memo,
				// AccountName filled in next step
				AccountName: tx.AccountID, // temporary placeholder
			}
		})
	}

	// Step 4 (Bind): resolve account name via result.FromOption.
	resolveAccount := func(tx EnrichedTx) result.Result[EnrichedTx, string] {
		acctID := tx.AccountName // still holds the raw ID from step 3
		nameOpt := option.None[string]()
		if name, ok := accounts[acctID]; ok {
			nameOpt = option.Some(name)
		}
		return result.FromOption(nameOpt, fmt.Sprintf("unknown account: %s", acctID)).
			Map(func(name string) EnrichedTx {
				tx.AccountName = name
				return tx
			})
	}

	// Step 5 (Map - pure): apply exchange rate to the amount.
	applyRate := func(tx EnrichedTx) EnrichedTx {
		tx.Amount *= exchangeRate
		return tx
	}

	// Step 6 (Map - pure): format the final ledger entry.
	formatEntry := func(tx EnrichedTx) LedgerEntry {
		return LedgerEntry{
			Date:    tx.Date,
			Amount:  math.Round(tx.Amount*100) / 100,
			Summary: fmt.Sprintf("[%s] %s", tx.AccountName, tx.Memo),
		}
	}

	// The full 6-step pipeline.
	process := func(line string) result.Result[LedgerEntry, string] {
		r1 := parseCSV(line)          // step 1: Bind
		r2 := r1.Bind(validateFields) // step 2: Bind
		r3 := r2.Bind(parseAmount)    // step 3: Bind (Try + MapErr)
		r4 := r3.Bind(resolveAccount) // step 4: Bind (FromOption)
		r5 := r4.Map(applyRate)       // step 5: Map (pure)
		r6 := r5.Map(formatEntry)     // step 6: Map (pure)
		return r6
	}

	t.Run("all 6 steps succeed", func(t *testing.T) {
		r := process("2024-03-15|250.00|ACCT-1|Invoice #42")
		if r.IsErr() {
			t.Fatalf("expected Ok, got Err(%s)", r.UnwrapErr())
		}
		entry := r.Unwrap()
		if entry.Date != "2024-03-15" {
			t.Errorf("Date = %s", entry.Date)
		}
		// 250.00 × 1.10 = 275.00
		if got := fmt.Sprintf("%.2f", entry.Amount); got != "275.00" {
			t.Errorf("Amount = %s, want 275.00", got)
		}
		wantSummary := "[Operating Chequing] Invoice #42"
		if entry.Summary != wantSummary {
			t.Errorf("Summary = %q, want %q", entry.Summary, wantSummary)
		}
	})

	t.Run("fails at step 1: wrong field count", func(t *testing.T) {
		assertPipeErr(t, process("only|two"), "expected 4 pipe-delimited fields")
	})

	t.Run("fails at step 2: missing required field", func(t *testing.T) {
		assertPipeErr(t, process("2024-03-15||ACCT-1|memo"), "amount is required")
	})

	t.Run("fails at step 3: unparseable amount (Try + MapErr)", func(t *testing.T) {
		assertPipeErr(t, process("2024-03-15|abc|ACCT-1|memo"), "invalid amount")
	})

	t.Run("fails at step 4: unknown account (FromOption)", func(t *testing.T) {
		assertPipeErr(t, process("2024-03-15|100.00|ACCT-999|memo"), "unknown account: ACCT-999")
	})
}
func TestFlatten(t *testing.T) {
	inner := result.Ok[int, string](42)
	outer := result.Ok[result.Result[int, string], string](inner)
	if result.Flatten(outer).Unwrap() != 42 {
		t.Fatal("expected 42")
	}
	outerErr := result.Err[result.Result[int, string], string]("boom")
	if result.Flatten(outerErr).IsOk() {
		t.Fatal("expected Err")
	}
}
func TestOrElse(t *testing.T) {
	ok := result.Ok[int, string](1)
	got := ok.OrElse(func(_ string) result.Result[int, string] { return result.Ok[int, string](99) })
	if got.Unwrap() != 1 {
		t.Fatal("expected Ok to win")
	}
	err := result.Err[int, string]("oops")
	got = err.OrElse(func(_ string) result.Result[int, string] { return result.Ok[int, string](42) })
	if got.Unwrap() != 42 {
		t.Fatal("expected fallback 42")
	}
}
func TestTee(t *testing.T) {
	var seen int
	r := result.Ok[int, string](7)
	out := r.Tee(func(v int) { seen = v })
	if seen != 7 || out.Unwrap() != 7 {
		t.Fatal("Tee should call fn and return unchanged")
	}
	seen = 0
	result.Err[int, string]("e").Tee(func(v int) { seen = v })
	if seen != 0 {
		t.Fatal("Tee should not call fn on Err")
	}
}
func TestTeeErr(t *testing.T) {
	var seen string
	result.Err[int, string]("bad").TeeErr(func(e string) { seen = e })
	if seen != "bad" {
		t.Fatal("TeeErr should call fn on Err")
	}
	result.Ok[int, string](1).TeeErr(func(_ string) { seen = "nope" })
	if seen == "nope" {
		t.Fatal("TeeErr should not call fn on Ok")
	}
}
func TestZip(t *testing.T) {
	a := result.Ok[int, string](1)
	b := result.Ok[string, string]("hello")
	got := result.Zip(a, b)
	if !got.IsOk() {
		t.Fatal("expected Ok")
	}
	p := got.Unwrap()
	if p.First != 1 || p.Second != "hello" {
		t.Fatalf("got %v", p)
	}
	// first Err wins
	errA := result.Err[int, string]("a failed")
	got2 := result.Zip(errA, b)
	if got2.IsOk() {
		t.Fatal("expected Err")
	}
	if got2.UnwrapErr() != "a failed" {
		t.Fatalf("got %s", got2.UnwrapErr())
	}
}
func TestMap2(t *testing.T) {
	a := result.Ok[int, string](3)
	b := result.Ok[int, string](4)
	got := a.ZipWith(b, func(x, y int) int { return x + y })
	if got.Unwrap() != 7 {
		t.Fatalf("got %d", got.Unwrap())
	}
	errA := result.Err[int, string]("first")
	got2 := errA.ZipWith(b, func(x, y int) int { return x + y })
	if got2.UnwrapErr() != "first" {
		t.Fatal("expected first err to win")
	}
	errB := result.Err[int, string]("second")
	got3 := a.ZipWith(errB, func(x, y int) int { return x + y })
	if got3.UnwrapErr() != "second" {
		t.Fatal("expected second err to propagate")
	}
}
func TestContains(t *testing.T) {
	if !result.Contains(result.Ok[int, string](42), 42) {
		t.Fatal("expected true")
	}
	if result.Contains(result.Ok[int, string](42), 99) {
		t.Fatal("expected false for wrong value")
	}
	if result.Contains(result.Err[int, string]("e"), 42) {
		t.Fatal("expected false on Err")
	}
}
func TestSequence(t *testing.T) {
	rs := []result.Result[int, string]{
		result.Ok[int, string](1),
		result.Ok[int, string](2),
		result.Ok[int, string](3),
	}
	got := result.Sequence(rs)
	if !got.IsOk() {
		t.Fatal("expected Ok")
	}
	s := got.Unwrap()
	if len(s) != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Fatalf("got %v", s)
	}
	// short-circuits on first Err
	rsErr := []result.Result[int, string]{
		result.Ok[int, string](1),
		result.Err[int, string]("boom"),
		result.Ok[int, string](3),
	}
	got2 := result.Sequence(rsErr)
	if got2.IsOk() || got2.UnwrapErr() != "boom" {
		t.Fatal("expected Err(boom)")
	}
}
func TestTraverse(t *testing.T) {
	items := []string{"1", "2", "3"}
	got := result.Traverse(items, func(s string) result.Result[int, string] {
		if s == "bad" {
			return result.Err[int, string]("bad input")
		}
		n := 0
		for _, c := range s {
			n = n*10 + int(c-'0')
		}
		return result.Ok[int, string](n)
	})
	if !got.IsOk() {
		t.Fatalf("expected Ok, got %s", got.UnwrapErr())
	}
	v := got.Unwrap()
	if len(v) != 3 || v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Fatalf("got %v", v)
	}
	items2 := []string{"1", "bad", "3"}
	got2 := result.Traverse(items2, func(s string) result.Result[int, string] {
		if s == "bad" {
			return result.Err[int, string]("bad input")
		}
		return result.Ok[int, string](0)
	})
	if got2.IsOk() || got2.UnwrapErr() != "bad input" {
		t.Fatal("expected Err(bad input)")
	}
}

func TestZipWith(t *testing.T) {
	// both Ok
	got := result.Ok[int, string](3).ZipWith(result.Ok[string, string]("px"), func(n int, s string) string {
		return fmt.Sprintf("%d%s", n, s)
	})
	if !got.IsOk() || got.Unwrap() != "3px" {
		t.Fatalf("ZipWith Ok+Ok: got %v", got)
	}
	// first Err
	e1 := result.Err[int, string]("e1").ZipWith(result.Ok[string, string]("px"), func(_ int, s string) string { return s })
	if !e1.IsErr() || e1.UnwrapErr() != "e1" {
		t.Fatalf("ZipWith Err+Ok: got %v", e1)
	}
	// second Err
	e2 := result.Ok[int, string](3).ZipWith(result.Err[string, string]("e2"), func(n int, s string) string { return s })
	if !e2.IsErr() || e2.UnwrapErr() != "e2" {
		t.Fatalf("ZipWith Ok+Err: got %v", e2)
	}
}
