# inquiry
1. endpoint : /api/iak/inquiry
2. request : models.RequestInquiry, 
3. jika ProductCategoryID = 4 (pln) dan type product id = 1(prepaid) maka 
    a. url dev : https://prepaid.iak.dev
        url prod : https://prepaid.iak.id
    b. endpoint : /api/inquiry-pln
    c. username : 082137789378
    d. assign data request ke models.ReqPaymentPrepaidIak
    e. api key dev : 51562ac44252c544AhiD
        api key prod : 78362ac44e3786a5QWrQ
    f. sign: md5({username}+{api_key}+{additional})
    g. additional = CustomerID
    h. method : post
    i. response ditampung pada models.RespPaymentPrepaidIAK, jika Data.SubscriberID == "" tampung pada models.RespWorkerIakUndefined, jika ResponseCode == "" tampung pada models.RespWorkerIakUndefinedI 
    j. tampung data yang "" tadi pada data yang menjadi acuan di models.RespPaymentPrepaidIAK dengan menempatkan data Data.Rc dan Data.Message pada RespPaymentPrepaidIAK.Rc dan RespPaymentPrepaidIAK.Message.
4. jika ProductCategoryID = 4 (pln) dan type product id = 2(postpaid) maka 
    a. url dev : https://testpostpaid.mobilepulsa.net
        url prod : https://mobilepulsa.net
    b. endpoint : /api/v1/bill/check
    c. username : 082137789378
    d. assign data request ke models.ReqInquiryPostpaidIAK
    e. Commands= "inq-pasca"
    f. api key dev : 51562ac44252c544AhiD
        api key prod : 78362ac44e3786a5QWrQ
    f. sign: md5({username}+{api_key}+{additional})
    g. additional = ReferenceNumber
    h. method : post
    i. response ditampung pada models.RespPaymentPrepaidIAK, jika Data.RefID == "" tampung pada models.RespWorkerIakUndefined, jika ResponseCode == "" tampung pada models.RespWorkerIakUndefinedI 
    j. jika Data.ResponseCode == "" tampung pada models.RespWorkerIakUndefinedII ,tampung data yang "" tadi pada data yang menjadi acuan di models.RespPaymentPostpaidIAK dengan menempatkan data Data.Rc dan Data.Message pada RespPaymentPostpaidIAK.Rc dan RespPaymentPostpaidIAK.Message.
5. konversi response menggunakan helpers.IAKConverterResponse
6. return data baku dengan format models.InquirytResult
