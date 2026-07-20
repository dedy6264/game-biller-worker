# inquiry
1. endpoint : /api/iak/payment
2. request : models.RequestPayment, 
3. url dev : https://prepaid.iak.dev
    url prod : https://prepaid.iak.id
4. enpoint : /api/top-up
5. username : 082137789378
6. assign data request ke models.ReqPaymentPrepaidIak
7. api key dev : 51562ac44252c544AhiD
    api key prod : 78362ac44e3786a5QWrQ
8. sign: md5({username}+{api_key}+{additional})
9. additional = reference number
10. method : post
11. response ditampung pada models.RespPaymentPrepaidIAK, jika Data.RefID == "" tampung pada models.RespWorkerIakUndefined, jika ResponseCode == "" tampung pada models.RespWorkerIakUndefinedI dan jika Data.ResponseCode="" tampung pada models.RespWorkerIakUndefinedII
12. tampung data yang "" tadi pada data yang menjadi acuan di models.RespPaymentPrepaidIAK dengan menempatkan data Data.Rc dan Data.Message pada RespPaymentPrepaidIAK.Rc dan RespPaymentPrepaidIAK.Message.
13. konversi response menggunakan helpers.IAKConverterResponse
14. return data baku dengan format models.PaymentResult
