$(document).ready(function() {

    $("#add-car").on("click", function() {
        addCar();
    });

    $("#update-car").on("click", function() {
        updateCar();
    });

    $("#delete-car").on("click", function() {
        deleteCar();
    });
});

// Add
function addCar() {
    const data = $("#product-form").serialize();

    $.ajax({
        type: "POST",
        url: "/add_car",
        data: data,
        success: function(response) {
            loadProductList();
            $("#product-form")[0].reset();
        },
        error: function(xhr) {
            alert("Error adding car: " + xhr.responseText);
        }
    });
}

// Update
function updateCar() {
    const data = $("#product-form").serialize();

    $.ajax({
        type: "POST",
        url: "/update_car",
        data: data,
        success: function(response) {
            loadProductList();
            $("#product-form")[0].reset();
        },
        error: function(xhr) {
            alert("Error updating car: " + xhr.responseText);
        }
    });
}

// Delete
function deleteCar() {
    const carId = $("#product-form").find("input[name='id']").val();

    if (!carId) {
        alert("Please select a car to delete.");
        return;
    }

    if (confirm("Are you sure you want to delete this car?")) {
        $.ajax({
            type: "POST",
            url: "/delete_car",
            data: { id: carId },
            success: function(response) {
                alert("Car deleted successfully");
                loadProductList();
                $("#product-form")[0].reset();
            },
            error: function(xhr) {
                alert("Error deleting car: " + xhr.responseText);
            }
        });
    }
}

// refreshes list without page reload

}
