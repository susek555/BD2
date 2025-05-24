// SaleOffer

export interface BaseOffer {
  id: string;
  // image: type to  be determined;
  name: string; // producer and model
  productionYear: number;
  mileage: number;
  color: string;
  price: number;
  isAuction: boolean;
}

export interface SaleOffer extends BaseOffer {
  isFavorite: boolean;
}

export interface HistoryOffer extends BaseOffer {
  dateEnd: string;
}
