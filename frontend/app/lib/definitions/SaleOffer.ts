// SaleOffer

export interface BaseOffer {
  id: string;
  // image: type to  be determined;
  name: string; // producer and model
  production_year: number;
  mileage: number;
  color: string;
  price: number;
  is_auction: boolean;
  main_url: string;
  can_modify: boolean;
}

export interface SaleOffer extends BaseOffer {
  is_liked: boolean;
}

export interface HistoryOffer extends BaseOffer {
  dateEnd: string;
  sellerRating?: number;
  sellerId: number;
  sellerName: string;
}
