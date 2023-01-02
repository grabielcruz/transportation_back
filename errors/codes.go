package errors_handler

// Database
const DB001 = "sql: no rows in result set"
const DB002 = "Could not begin transaction"
const DB003 = "Could not commit transaction"

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

// Currencies
const CU001 = "Could not delete VED or USD currency"
const CU002 = "Currency code should be 3 upper case letters"
const CU003 = "Currency already exists"
const CU004 = "Currency is being used"

// Transactions
const TR001 = "Could not get balance from account"
const TR002 = "Transaction should not generate a negative balance"
const TR003 = "The transaction requested is not the last transaction"
const TR004 = "No transaction found in database"
const TR005 = "Could not update account's balance"
const TR006 = "Could not insert record into trashed transactions table"
const TR007 = "Could not get transactions"
const TR008 = "Could not get trashed transactions"
const TR009 = "Could not read transaction"
const TR010 = "Could not read trashed transaction"
const TR011 = "Could not delete trashed transaction"
const TR012 = "Could not create restored transaction"
