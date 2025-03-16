// src/components/AddItem.tsx
import React, { useState } from "react";

interface AddItemProps {
  onAdd: (item: { ID: string; Quantity: number }) => void;
}

const AddItem: React.FC<AddItemProps> = ({ onAdd }) => {
  const [selectedItem, setSelectedItem] = useState("item-1");
  const [quantity, setQuantity] = useState(1);

  const handleAdd = () => {
    onAdd({ ID: selectedItem, Quantity: quantity });
  };

  return (
    <div className="mb-4 flex flex-col rounded-lg bg-white p-4 shadow">
      <label className="mb-2 font-medium">选择商品:</label>
      <select
        className="mb-2 rounded border p-2"
        value={selectedItem}
        onChange={(e) => setSelectedItem(e.target.value)}
      >
        <option value="item-1">item-1</option>
        <option value="item-2">item-2</option>
        <option value="item-3">item-3</option>
      </select>
      <label className="mb-2 font-medium">数量:</label>
      <input
        type="number"
        className="mb-2 rounded border p-2"
        value={quantity}
        onChange={(e) => setQuantity(Number(e.target.value))}
      />
      <button
        className="rounded bg-blue-500 p-2 text-white hover:bg-blue-600"
        onClick={handleAdd}
      >
        添加商品
      </button>
    </div>
  );
};

export default AddItem;
