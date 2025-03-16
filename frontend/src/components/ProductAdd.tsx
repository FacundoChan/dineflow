import React, { useState } from "react";

interface ProductAddProps {
  onAdd: (item: { id: string; quantity: number }) => void;
}

const ProductAdd: React.FC<ProductAddProps> = ({ onAdd }) => {
  const [selectedItem, setSelectedItem] = useState("item-1");
  const [quantity, setQuantity] = useState(1);

  const handleAdd = () => {
    if (quantity > 0) {
      onAdd({ id: selectedItem, quantity });
      setQuantity(1);
    }
  };

  return (
    <div className="mb-4 rounded border p-4 shadow">
      <select
        value={selectedItem}
        onChange={(e) => setSelectedItem(e.target.value)}
        className="mr-2 rounded border p-2"
      >
        <option value="item-1">item-1</option>
        <option value="item-2">item-2</option>
        <option value="item-3">item-3</option>
      </select>
      <input
        type="number"
        value={quantity}
        onChange={(e) => setQuantity(parseInt(e.target.value))}
        className="mr-2 w-20 rounded border p-2"
        min="1"
      />
      <button
        onClick={handleAdd}
        className="rounded bg-blue-500 p-2 text-white"
      >
        添加商品
      </button>
    </div>
  );
};

export default ProductAdd;
