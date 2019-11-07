
// Click on pencil to make title editable
$("#pencil").click(function(event){
    event.preventDefault()
    $("#title-editable").show()
    $("#title-fix").hide()
})

// Enable submit button if anything changes
$(".updatable").on("keyup change", function() {
    console.log("Hello from updatable")
    $("#submit").removeAttr("disabled")
})