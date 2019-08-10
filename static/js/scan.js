
// https://www.dynamsoft.com/CustomerPortal/Portal/TrialLicense.aspx
BarcodeReader.licenseKey = 't0068NQAAAG6yi5kCx2bAaEALB36Yhsk9OgzLadNZjJXj8OgZifmUZBZvO9oGvPMAHvNjOSwJIdhaJZcBrTQL5usgBfC1trE=';
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

