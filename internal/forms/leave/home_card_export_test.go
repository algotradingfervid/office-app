package leave

// BalanceCardAt lets the external tests render the card with a fixed clock
// (they cannot be in package leave: testapp imports this module).
var BalanceCardAt = balanceCard
