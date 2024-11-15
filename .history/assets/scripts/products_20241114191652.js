console.log("JS loaded");
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
    console.log("Data being sent:", data);

    $.ajax({
        type: "POST",
        url: "/add_car",
        data: data,
        success: function(response) {
            console.log("Server response:", response);
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
    const carId = $("#car-id").val();
    
    // Collect form data, add inactive explicitly if unchecked
    const data = $("#product-form").serializeArray();
    if (!$("#inactive").is(":checked")) {
        data.push({ name: "inactive", value: "0" });
    }
    
    if (!carId) {
        alert("Please select a car to update.");
        return;
    }

    $.ajax({
        type: "PUT",
        url: `/update_car/${carId}`,
        data: $.param(data), // Convert array to URL-encoded string
        success: function(response) {
            console.log("Car updated successfully:", response);
            loadProductList(); // Reload product list to reflect changes
            $("#product-form")[0].reset(); // Clear the form
        },
        error: function(xhr) {
            alert("Error updating car: " + xhr.responseText);
        }
    });
}

function selectCar(id, name, image, quantity, price, inactive) {
    console.log("Selected Car Details:", { id, name, image, quantity, price, inactive });

    // Populate the form with selected car's data
    $("#car-id").val(id);
    $("#name").val(name);
    $("#image").val(image);
    $("#quantity").val(quantity);
    $("#price").val(price);
    $("#inactive").prop("checked", inactive === 1);

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
    console.log("loadProductList called");
    $.ajax({
        url: "/get_all_products",
        method: "GET",
        success: function(response) {
            console.log("Received products:", response);
            const tableBody = $("#product-list tbody");
            tableBody.empty(); // Clear existing rows
            response.products.forEach(product => {
                const inactiveStatus = product.Inactive === 1 ? "Yes" : "No";
                const row = $(`
                    <tr>
                        <td>${product.Name}</td>
                        <td>${product.Image}</td>
                        <td>${product.QuantityInStock}</td>
                        <td>${product.Price}</td>
                        <td>${inactiveStatus}</td>
                    </tr>
                `);
                
                // Attach the click event to each row
                row.on("click", function() {
                    selectCar(product.ID, product.Name, product.Image, product.QuantityInStock, product.Price, product.Inactive);
                });
                
                tableBody.append(row);
            });
        },
        error: function(error) {
            console.error("Error loading product list:", error);
        }
    });
}

