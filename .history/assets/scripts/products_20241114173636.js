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
function loadProductList() {
    $.ajax({
        url: "/get_all_products",  // Ensure this endpoint returns a JSON of all products
        method: "GET",
        success: function(response) {
            const tableBody = $("#product-list tbody");
            tableBody.empty(); // Clear existing rows
            response.products.forEach(product => {
                // Assuming the backend returns an array of product objects
                const inactiveStatus = product.inactive === 1 ? "Yes" : "No";
                tableBody.append(`
                    <tr>
                        <td>${product.name}</td>
                        <td>${product.image}</td>
                        <td>${product.quantity}</td>
                        <td>${product.price}</td>
                        <td>${inactiveStatus}</td>
                    </tr>
                `);
            });
        },
        error: function(error) {
            console.error("Error loading product list:", error);
        }
    });
}

