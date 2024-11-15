$(document).ready(function() {
    // Helper function to check required fields
    function checkRequiredFields() {
        const name = $("#name").val().trim();
        const image = $("#image").val().trim();
        const price = $("#price").val().trim();
        
        if (!name) {
            alert("Product name must not be blank.");
            $("#name").focus();
            return false;
        }
        if (!image) {
            alert("Image name must not be blank.");
            $("#image").focus();
            return false;
        }
        if (!price) {
            alert("Price must not be blank.");
            $("#price").focus();
            return false;
        }
        return true;
    }

    // Add Product
    $("#add-product").click(function() {
        if (!checkRequiredFields()) return;
        
        const productData = {
            name: $("#name").val(),
            image: $("#image").val(),
            quantity: $("#quantity").val(),
            price: $("#price").val(),
            inactive: $("#inactive").is(":checked") ? 1 : 0
        };

        $.ajax({
            type: "POST",
            url: "/add_product",
            data: JSON.stringify(productData),
            contentType: "application/json",
            success: function(response) {
                alert("Product added successfully.");
                loadProductList();
                $("#product-form")[0].reset(); // Clear form
            },
            error: function(xhr) {
                alert("Error adding product.");
            }
        });
    });

    // Update Product
    $("#update-product").click(function() {
        if (!checkRequiredFields()) return;

        const productData = {
            id: selectedProductID,
            name: $("#name").val(),
            image: $("#image").val(),
            quantity: $("#quantity").val(),
            price: $("#price").val(),
            inactive: $("#inactive").is(":checked") ? 1 : 0
        };

        $.ajax({
            type: "PUT",
            url: "/update_product/" + selectedProductID,
            data: JSON.stringify(productData),
            contentType: "application/json",
            success: function(response) {
                alert("Product updated successfully.");
                loadProductList();
                $("#product-form")[0].reset(); // Clear form
            },
            error: function(xhr) {
                alert("Error updating product.");
            }
        });
    });

    // Delete Product
    $("#delete-product").click(function() {
        if (!selectedProductID) {
            alert("Please select a product to delete.");
            return;
        }

        if (!confirm("Are you sure you want to delete this product?")) return;

        $.ajax({
            type: "DELETE",
            url: "/delete_product/" + selectedProductID,
            success: function(response) {
                alert("Product deleted successfully.");
                loadProductList();
                $("#product-form")[0].reset(); // Clear form
            },
            error: function(xhr) {
                alert("Error deleting product.");
            }
        });
    });

    // Load Product List
    function loadProductList() {
        $.ajax({
            type: "GET",
            url: "/products",
            success: function(products) {
                let productRows = "";
                products.forEach(product => {
                    productRows += `
                        <tr onclick="selectProduct(${product.id}, '${product.name}', '${product.image}', ${product.quantity}, ${product.price}, ${product.inactive})">
                            <td>${product.name}</td>
                            <td>${product.image}</td>
                            <td>${product.quantity}</td>
                            <td>${product.price.toFixed(2)}</td>
                            <td>${product.inactive ? 'Yes' : 'No'}</td>
                        </tr>
                    `;
                });
                $("#product-list tbody").html(productRows);
            },
            error: function(xhr) {
                alert("Error loading products.");
            }
        });
    }

    // Populate form when product row is clicked
    window.selectProduct = function(id, name, image, quantity, price, inactive) {
        selectedProductID = id;
        $("#name").val(name);
        $("#image").val(image);
        $("#quantity").val(quantity);
        $("#price").val(price.toFixed(2));
        $("#inactive").prop("checked", inactive === 1);
    };

    // Initial load
    loadProductList();
});
