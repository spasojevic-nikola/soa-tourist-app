export interface OrderItem {
    tourId: string;
    name: string;
    price: number;
  }
  
 
  export interface ShoppingCart {
    id: string;
    userId: number;
    items: OrderItem[];
    total: number;
    updatedAt: string;
  }