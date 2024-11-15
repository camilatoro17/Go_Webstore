$(document).ready(function() {
    // Load products on page load
    loadProductList();

    // Attach event listeners
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

// Add Car
function addCar() {
    const data = $("#product-form").serialize();

    $.ajax({
        type: "POST",
        url: "/add_car",
        data: data,
        success: function(response) {
            loadProductList(); // Reload product list after adding
            $("#product-form")[0].reset(); // Reset the form
        },
        error: function(xhr) {
            alert("Error adding car: " + xhr.responseText);
        }
    });
}

// Update Car
function updateCar() {
    const carId = $("#car-id").val(); // Get the ID from the hidden input
    const data = $("#product-form").serialize();

    $.ajax({
        type: "PUT",
        url: `/update_car/${carId}`, // Include the car ID in the URL
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


// Delete Car
function deleteCar() {
    const carId = $("#product-form").find("input[name='id']").val();

    if (!carId) {
        alert("Please select a car to delete.");
        return;
    }

    if (confirm("Are you sure you want to delete this car?")) {
        $.ajax({
            type: "DELETE",
            url: "/delete_car/" + carId,
            success: function(response) {
                loadProductList(); // Reload product list after deleting
                $("#product-form")[0].reset(); // Reset the form
            },
            error: function(xhr) {
                alert("Error deleting car: " + xhr.responseText);
            }
        });
    }
}

// Load Product List
function loadProductList() {
    $.ajax({
        url: "/get_all_products",
        method: "GET",
        success: function(response) {
            const tableBody = $("#product-list tbody");
            tableBody.empty(); // Clear existing rows
            response.products.forEach(product => {
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
