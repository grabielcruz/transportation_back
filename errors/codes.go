package errors_handler

// Database
// const DB001 = "sql: no rows in result set"
const DB001 = "Record not found in database"
const DB002 = "Could not begin transaction"
const DB003 = "Could not commit transaction"
const DB004 = "Could not count records"
const DB005 = "Could not get records"
const DB007 = "Could not insert record"
const DB009 = "Could not update record"

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

// Persons
const PE001 = "Document already in use"
const PE002 = "Person does not exists"

// Currencies
const CU001 = "Could not delete VED or USD currency"
const CU002 = "Currency code should be 3 upper case letters"
const CU003 = "Currency already exists"
const CU004 = "Currency is being used"

// Foreign Key error
const CU005 = "Currency it is not registered in database"

// Transactions
const TR001 = "Could not get balance from account"
const TR002 = "Transaction should not generate a negative balance"
const TR003 = "The transaction requested is not the last transaction"

// const TR004 = "No transaction found in database"
const TR005 = "Could not update account's balance"
const TR006 = "New balance and updated balance missmatch, oldBalance = %v, newBalance = %v, updatedBalance = %v"
const TR007 = "Transaction should have a person"
const TR008 = "Transaction should have an amount different from zero"
const TR009 = "Fee should be between 0 and 1"
const TR010 = "Person's account does not belong to the person specified"
const TR011 = "Currency's mismatch"
const TR012 = "Money account does not exist"
const TR013 = "Pending bill does not exist"
const TR014 = "Positive transactions should not have a fee"

// Bills
const BL001 = "Could not request empty set of bills"
const BL002 = "Can not create bill with amount of zero"
const BL003 = "Can not delete pending bill associated to transaction"

// const BL003 = "Person with the specified uuid does not exists"
// const BL004 = "Currency it is not registered in database"

//Person Accounts
const PA001 = "Person with the specified uuid does not exists"
const PA002 = "Person account does not exist"
