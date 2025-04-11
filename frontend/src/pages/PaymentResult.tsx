import React, { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import apiClient from "@/services/api-client";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

interface OrderItem {
  id: string;
  quantity: number;
  name?: string;
  price?: number;
  img_urls?: string[];
}

interface OrderData {
  id: string;
  status: string;
  items: OrderItem[];
}

interface ProductMap {
  [id: string]: {
    name: string;
    price: number;
    img_urls: string[];
  };
}

const PaymentResult: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<string>("");
  const [hasAlerted, setHasAlerted] = useState(false);
  const [order, setOrder] = useState<OrderData | null>(null);
  const [products, setProducts] = useState<ProductMap>({});
  const customerID = searchParams.get("customerID");
  const orderID = searchParams.get("orderID");

  useEffect(() => {
    const fetchOrderAndProducts = async () => {
      if (!customerID || !orderID) return;

      try {
        const [orderRes, productsRes] = await Promise.all([
          apiClient.get(`/api/customer/${customerID}/orders/${orderID}`),
          apiClient.get("/api/products"),
        ]);

        const orderData = orderRes.data.data.Order;
        const productList = productsRes.data.data.products;
        const productMap: ProductMap = {};

        for (const product of productList) {
          productMap[product.id] = {
            name: product.name,
            price: product.price,
            img_urls: product.img_urls || [],
          };
        }

        const itemsWithDetails = orderData.items.map((item: OrderItem) => ({
          ...item,
          name: productMap[item.id]?.name || item.id,
          price: productMap[item.id]?.price || 0,
          img_urls: productMap[item.id]?.img_urls || [],
        }));

        setOrder({
          id: orderData.id,
          status: orderData.status,
          items: itemsWithDetails,
        });
        setStatus(orderData.status);
        setProducts(productMap);
      } catch (error) {
        console.error("获取订单或商品信息失败", error);
        setStatus("error");
      }
    };

    fetchOrderAndProducts();
  }, [customerID, orderID]);

  useEffect(() => {
    if (status !== "paid") return;

    const interval = setInterval(async () => {
      try {
        const response = await apiClient.get(
          `/api/customer/${customerID}/orders/${orderID}`,
        );
        const order = response.data.data.Order;
        if (order.status === "ready" && !hasAlerted) {
          setHasAlerted(true);
          setStatus("ready");
        }
      } catch (error) {
        console.error("轮询失败", error);
      }
    }, 5000);

    return () => clearInterval(interval);
  }, [status, customerID, orderID, hasAlerted]);

  useEffect(() => {
    if (!status) return;

    switch (status) {
      case "paid":
        toast.success("✅ 支付成功", {
          description: "我们已收到您的订单",
        });
        break;
      case "waiting_for_payment":
        toast.info("⌛ 等待支付中", {
          description: "请在 5 分钟内完成付款",
        });
        break;
      case "ready":
        toast.success("🎉 订单已完成", {
          description: "感谢您的购买！",
          className: "text-center",
        });
        break;
      case "error":
        toast.error("❌ 查询失败", {
          description: "请稍后再试",
        });
        break;
    }
  }, [status]);

  const getStatusBadge = () => {
    switch (status) {
      case "paid":
        return <Badge className="bg-green-500">✅ 支付成功</Badge>;
      case "waiting_for_payment":
        return <Badge className="bg-yellow-500">⌛ 等待支付中</Badge>;
      case "ready":
        return <Badge className="bg-blue-500">🎉 订单已完成</Badge>;
      case "error":
        return <Badge className="bg-red-500">❌ 查询失败</Badge>;
      default:
        return <Badge>🔄 查询中...</Badge>;
    }
  };

  return (
    <div className="container mx-auto p-4">
      <div className="mb-4 text-center">{getStatusBadge()}</div>

      {!order ? (
        <Skeleton className="mx-auto h-24 w-full max-w-2xl" />
      ) : (
        <div className="space-y-4">
          {order.items.map((item) => {
            const subtotal = item.quantity * (item.price || 0);
            return (
              <Card
                key={item.id}
                className="mx-auto flex w-full max-w-2xl flex-row items-center justify-between px-4"
              >
                <div className="flex items-center gap-3">
                  {item.img_urls?.[0] && (
                    <img
                      src={item.img_urls[0]}
                      alt={item.name}
                      className="h-12 w-12 rounded-md object-cover"
                    />
                  )}
                  <div>
                    <div className="text-lg font-semibold">{item.name}</div>
                    <div className="text-sm text-gray-500">
                      数量：{item.quantity}
                    </div>
                  </div>
                </div>
                <div className="text-right whitespace-nowrap">
                  <div className="text-lg font-bold">
                    HK${subtotal.toFixed(2)}
                  </div>
                  <div className="text-muted-foreground text-sm">
                    HK${(item.price || 0).toFixed(2)} 每个
                  </div>
                </div>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default PaymentResult;
