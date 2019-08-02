
// https://www.dynamsoft.com/CustomerPortal/Portal/TrialLicense.aspx
BarcodeReader.licenseKey = 't0068NQAAAA9xqcBsrb1DuHpiuyXSc3C7EdWbaxEagkUKVVNMdGTtUkeW4rNjxGV0DkaMn5OldSQJsnhkVKrT1os5WZDlvIU=';
let scanner = new BarcodeReader.Scanner({
    onFrameRead: results => {console.log(results);},
    onNewCodeRead: (txt, result) => {window.location.href = "/product?gtin=" + txt;}
});
scanner.open().catch(ex=>{
    console.log(ex);
    alert(ex.message || ex);
    scanner.close();
});