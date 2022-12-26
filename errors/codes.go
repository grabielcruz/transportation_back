package errors_handler

// Database
const DB001 = "sql: no rows in result set"

// Reading error
const RE001 = "Unable to read body of the request"

// Unmarshal error
const UM001 = "Invalid data type"

// Validation error
// Validation error works with a custom message
const VA001 = "Validation error"

// Invalid UUID error
// Invalid UUID error works with a custom message
const UI001 = "Invalid UUID"

// Service error
const SE001 = "Service error"

// Querystring error
const QS001 = "Query string error"

// Transactions
const TR001 = "Could not get balance from account"
const TR002 = "Transaction should not generate a negative balance"
const TR003 = "The transaction requested is not the last transaction"
const TR004 = "No transaction found in database"
