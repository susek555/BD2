import { Review } from '../definitions';

export async function fetchReviewsByReviewer(id: number): Promise<Review[]> {
  const url = `${process.env.URL}/api/reviews/reviewer/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch history offers');
  }

  return response.json();
}

export async function fetchReviewsByReviewee(id: number): Promise<Review[]> {
  const url = `${process.env.URL}/api/reviews/reviewee/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch history offers');
  }

  return response.json();
}
