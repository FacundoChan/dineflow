import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

type Product = {
  id: string;
  name: string;
  quantity: number;
  price: number;
  img_urls?: string[];
};

type SelectedItem = {
  id: string;
  quantity: number;
};

interface ProductListProps {
  onChange: (items: SelectedItem[]) => void;
  products: Product[];
  setProducts: React.Dispatch<React.SetStateAction<Product[]>>;
}

const ProductList: React.FC<ProductListProps> = ({
  onChange,
  products,
  setProducts,
}) => {
  const [selected, setSelected] = useState<Record<string, number>>({});

  useEffect(() => {
    axios
      .get("http://127.0.0.1:8282/api/products")
      .then((res) => {
        setProducts(res.data?.data?.products || []);
      })
      .catch(console.error);
  }, []);

  const updateQuantity = (id: string, delta: number, stock: number) => {
    setSelected((prev) => {
      const newQty = Math.min(stock, Math.max(0, (prev[id] || 0) + delta));
      const updated = { ...prev, [id]: newQty };
      onChange(
        Object.entries(updated)
          .filter(([, quantity]) => quantity > 0)
          .map(([id, quantity]) => ({ id, quantity })),
      );
      return updated;
    });
  };

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4">
      {products.map((product) => (
        <Card key={product.id}>
          <CardHeader>
            <CardTitle className="text-xl font-bold">{product.name}</CardTitle>
            {product.img_urls?.[0] && (
              <img
                src={product.img_urls[0]}
                className="mx-auto mb-2 flex h-45 rounded-xl object-cover"
              />
            )}
          </CardHeader>
          <CardContent>
            <p>库存: {product.quantity}</p>
            <p className="text-muted-foreground">
              单价: HK${product.price.toFixed(2)}
            </p>
            <div className="mt-2 flex items-center justify-center space-x-2">
              <Button
                onClick={() => updateQuantity(product.id, -1, product.quantity)}
              >
                -
              </Button>
              <span>{selected[product.id] || 0}</span>
              <Button
                onClick={() => updateQuantity(product.id, 1, product.quantity)}
              >
                +
              </Button>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
};

export default ProductList;
