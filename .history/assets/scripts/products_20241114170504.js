$(document).ready(function() {
    // Add Car
    $("#add-product").click(function() {
        const productData = {
            name: $("#name").val(),
            image: $("#image").val(),
            quantity: $("#quantity").val(),
            price: $("#price").val(),
            inactive: $("#inactive").is(":checked") ? 1 : 0
        };

        $.ajax({
            type: "POST",
            url: "/add_car",
            data: JSON.stringify(productData),
            contentType: "application/json",
            success: function(response) {
                alert(response.message);
                loadProductList();
                clearForm();
            },
            error: function(xhr) {
                alert("Error adding car: " + xhr.responseText);
            }
        });
    });

    // Update Car
    $("#update-product").click(function() {
        const productId = $("#product-id").val();
        const productData = {
            name: $("#name").val(),
            image: $("#image").val(),
            quantity: $("#quantity").val(),
            price: $("#price").val(),
            inactive: $("#inactive").is(":checked") ? 1 : 0
        };

        $.ajax({
            type: "PUT",
            url: `/update_car/${productId}`,
            data: JSON.stringify(productData),
            contentType: "application/json",
            success: function(response) {
                alert(response.message);
                loadProductList();
                clearForm();
            },
            error: function(xhr) {
                alert("Error updating car: " + xhr.responseText);
            }
        });
    });

    // Delete Car
    $("#delete-product").click(function() {
        const productId = $("#product-id").val();

        if (confirm("Are you sure you want to delete this car?")) {
            $.ajax({
                type: "DELETE",
                url: `/delete_car/${productId}`,
                success: function(response) {
                    alert(response.message);
                    loadProductList();
                    clearForm();
                },
                error: function(xhr) {
                    alert("Error deleting car: " + xhr.responseText);
                }
            });
        }
    });

    // load without refreshing
    function loadProductList() {
        
    }

    // Function to clear form fields
    function clearForm() {
        $("#name").val("");
        $("#image").val("");
        $("#quantity").val("");
        $("#price").val("");
        $("#inactive").prop("checked", false);
        $("#product-id").val(""); 
    }
});
