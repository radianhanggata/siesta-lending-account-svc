package lending

import "github.com/radianhanggata/siesta-coding-test/lending-account-svc/internalerror"

var ConfigNotFound = &internalerror.Response{Code: 404, Message: "config not found"}
var RepaymentNotFound = &internalerror.Response{Code: 404, Message: "repayment not found"}
