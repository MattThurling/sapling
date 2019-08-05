
// https://www.dynamsoft.com/CustomerPortal/Portal/TrialLicense.aspx
BarcodeReader.licenseKey = 't0068NQAAAA9xqcBsrb1DuHpiuyXSc3C7EdWbaxEagkUKVVNMdGTtUkeW4rNjxGV0DkaMn5OldSQJsnhkVKrT1os5WZDlvIU=';
let scanner = new BarcodeReader.Scanner({
    onFrameRead: results => {
        hide();
        console.log(results);
    },
    onNewCodeRead: (txt, result) => {parent.window.location.href = "/product?gtin=" + txt;}
});
scanner.open().catch(ex=>{
    console.log(ex);
    alert(ex.message || ex);
    scanner.close();
});

function hide () {
    console.log("Hiding");
    document.getElementsByClassName("dbrScanner-sel-camera")[0].style.display = "none";
    document.getElementsByClassName("dbrScanner-sel-resolution")[0].style.display = "none";
}

