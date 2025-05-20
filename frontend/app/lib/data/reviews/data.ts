import { RatingPercentages, Review, ReviewSearchParams } from '../../definitions/reviews';

export async function fetchReviewsByReviewer(
  id: number,
  searchParams: ReviewSearchParams,
): Promise<Review[]> {
  const url = `${process.env.URL}/api/reviews/reviewer/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch history offers');
  }

  console.log('fetchReviewsByReviewer search params', searchParams);

  return response.json();
}

export async function fetchReviewsByReviewee(
  id: number,
  searchParams: ReviewSearchParams,
): Promise<Review[]> {
  const url = `${process.env.URL}/api/reviews/reviewee/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch history offers');
  }

  console.log('fetchReviewsByReviewee search params', searchParams);

  return response.json();
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
