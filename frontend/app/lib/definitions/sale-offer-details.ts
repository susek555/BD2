export type SaleOfferDetails = {
  id: number;
  name: string; // producer and model
  price?: number;
  isAuction: boolean;
  auctionData?: AuctionData;
  isActive: boolean;
  imagesURLs: string[]; // URLs to images
  details: OfferDetails[];
  description: string;
  // location: string;
  sellerName: string;
  sellerId: number;
  is_favourite: boolean;
  can_edit: boolean;
  can_delete: boolean;
};

export type AuctionData = {
  currentBid: number;
  endDate: Date;
};

export type OfferDetails = {
  name: string;
  value: string;
};
