// Reviews

interface ReviewUser {
  id: number;
  username: string;
}

interface BaseReview {
  description: string;
  rating: number;
}

export interface Review extends BaseReview {
  id: number;
  date: string;
  reviewer: ReviewUser;
  reviewee: ReviewUser;
}

export interface NewReview extends BaseReview {
  revieweeId: number;
}

export interface UpdatedReview extends BaseReview {
  id: number
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
