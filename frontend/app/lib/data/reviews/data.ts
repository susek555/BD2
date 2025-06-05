import {
  RatingPercentages,
  ReviewPage,
  ReviewSearchParams,
} from '../../definitions/reviews';

export async function fetchReviewsByReviewer(
  id: number,
  searchParams: ReviewSearchParams,
): Promise<ReviewPage> {
  const url = `${process.env.URL}/api/reviews/reviewer/${id}`;

  const response = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(searchParams),
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch reviews by ${id}`);
  }

  const data = await response.json();

  const mappedResponse: ReviewPage = {
    reviews: data.reviews.map((review: any) => ({
      id: review.id,
      rating: review.rating,
      description: review.description,
      date: review.review_date,
      reviewer: review.reviewer,
      reviewee: review.reviewee,
    })),
    pagination: {
      totalPages: data.pagination.total_pages,
      totalRecords: data.pagination.total_records,
    },
  };

  return mappedResponse;
}

export async function fetchReviewsByReviewee(
  id: number,
  searchParams: ReviewSearchParams,
): Promise<ReviewPage> {
  const url = `${process.env.URL}/api/reviews/reviewee/${id}`;

  const response = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(searchParams),
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch reviews about ${id}`);
  }

  const data = await response.json();

  const mappedResponse: ReviewPage = {
    reviews: data.reviews.map((review: any) => ({
      id: review.id,
      rating: review.rating,
      description: review.description,
      date: review.review_date,
      reviewer: review.reviewer,
      reviewee: review.reviewee,
    })),
    pagination: {
      totalPages: data.pagination.total_pages,
      totalRecords: data.pagination.total_records,
    },
  };

  return mappedResponse;
}

export async function fetchAverageRating(id: number): Promise<number> {
  const url = `${process.env.URL}/api/reviews/average/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch average rating');
  }

  return response.json();
}

export async function fetchRatingDistribution(
  id: number,
): Promise<RatingPercentages> {
  const url = `${process.env.URL}/api/reviews/distribution/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch rating distribution');
  }

  return response.json();
}
