    document.addEventListener('DOMContentLoaded', function() {
        loadItems();
        loadPurchaseHistory();

        document.getElementById('itemForm').addEventListener('submit', function(event) {
            event.preventDefault();
            const itemId = document.getElementById('itemId').value;
            const name = document.getElementById('name').value.trim();
            const price = parseFloat(document.getElementById('price').value);
            const quantity = parseInt(document.getElementById('quantity').value, 10);

            if (name && price >= 0 && quantity >= 0) {
                if (itemId) {
                    updateItem(itemId, { name, price, quantity });
                } else {
                    createItem({ name, price, quantity });
                }
            } else {
                alert('Please provide valid inputs.');
            }
        });

        document.getElementById('purchaseForm').addEventListener('submit', function(event) {
            event.preventDefault();
            const itemId = parseInt(document.getElementById('purchaseItemId').value, 10);
            const quantity = parseInt(document.getElementById('purchaseQuantity').value, 10);

            if (itemId >= 0 && quantity > 0) {
                purchaseItem({ item_id: itemId, quantity });
            } else {
                alert('Please provide valid inputs.');
            }
        });
    });

    function loadItems() {
        fetch('/items')
            .then(response => response.json())
            .then(data => {
                const itemsTableBody = document.getElementById('itemsTable').querySelector('tbody');
                itemsTableBody.innerHTML = '';
                data.items.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td class="px-6 py-4 border">${item.item_id}</td>
                            <td class="px-6 py-4 border">${item.name}</td>
                            <td class="px-6 py-4 border">${item.price.toFixed(2)}</td>
                            <td class="px-6 py-4 border">${item.quantity}</td>
                            <td class="px-6 py-4 border">
                                <button class="text-blue-500 hover:underline" onclick="editItem(${item.item_id})"><i class="fas fa-edit"></i> Edit</button>
                                <button class="text-red-500 hover:underline" onclick="deleteItem(${item.item_id})"><i class="fas fa-trash"></i> Delete</button>
                            </td>
                    `;
                    itemsTableBody.appendChild(row);
                });
            })
            .catch(() => alert('Error loading items.'));
    }

    function loadPurchaseHistory() {
    fetch('/items')
        .then(response => response.json())
        .then(data => {
            const purchaseHistoryTableBody = document.getElementById('purchaseHistoryTable').querySelector('tbody');
            purchaseHistoryTableBody.innerHTML = '';
            data.purchase_history.forEach(purchase => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td class="px-6 py-4 border">${purchase.item_id}</td>
                            <td class="px-6 py-4 border">${purchase.quantity}</td>
                            <td class="px-6 py-4 border">${purchase.total_price.toFixed(2)}</td>
                `;
                purchaseHistoryTableBody.appendChild(row);
            });
        });
}
    function createItem(item) {
        fetch('/items', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item),
        })
        .then(response => {
            if (response.ok) {
                loadItems();
                document.getElementById('itemForm').reset();
            } else {
                alert('Error creating item.');
            }
        });
    }

    function updateItem(id, item) {
        fetch(`/items/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item),
        })
        .then(response => {
            if (response.ok) {
                loadItems();
                document.getElementById('itemForm').reset();
            } else {
                alert('Error updating item.');
            }
        });
    }

    function deleteItem(id) {
        fetch(`/items/${id}`, { method: 'DELETE' })
        .then(response => {
            if (response.ok) {
                loadItems();
            } else {
                alert('Error deleting item.');
            }
        });
    }

    function editItem(id) {
        fetch(`/items/${id}`)
            .then(response => response.json())
            .then(item => {
                document.getElementById('itemId').value = item.item_id;
                document.getElementById('name').value = item.name;
                document.getElementById('price').value = item.price;
                document.getElementById('quantity').value = item.quantity;
            })
            .catch(() => alert('Error loading item.'));
    }

    function purchaseItem(purchase) {
        fetch('/purchase', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(purchase),
        })
        .then(response => {
            if (response.ok) {
                loadItems();
                loadPurchaseHistory();
                document.getElementById('purchaseForm').reset();
            } else {
                alert('Error purchasing item.');
            }
        });
    }