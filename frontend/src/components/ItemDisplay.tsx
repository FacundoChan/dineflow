// src/components/ItemDisplay.tsx
import React from "react";

interface ItemDisplayProps {
  item: { ID: string; Quantity: number };
}

const ItemDisplay: React.FC<ItemDisplayProps> = ({ item }) => {
  return (
    <div className="mb-2 flex items-center justify-between rounded bg-gray-100 p-4">
      <span className="font-medium">{item.ID}</span>
      <span>数量: {item.Quantity}</span>
    </div>
  );
};

export default ItemDisplay;
