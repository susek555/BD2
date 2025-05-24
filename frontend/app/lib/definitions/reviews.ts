// Reviews

interface ReviewUser {
  id: number;
  username: string;
}

export interface Review {
  id: number;
  description: string;
  rating: number;
  date: string;
  reviewer: ReviewUser;
  reviewee: ReviewUser;
}

export interface ReviewSearchParams {
  is_order_desc: boolean;
  order_key: 'rating' | 'date';
  pagination: {
    page: number;
    page_size: number;
  };
  ratings?: number[];
}

export interface RatingPercentages {
  '1': number;
  '2': number;
  '3': number;
  '4': number;
  '5': number;
}
