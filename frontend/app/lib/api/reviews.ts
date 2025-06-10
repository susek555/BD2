import { NewReview, Review, UpdatedReview } from '../definitions/reviews';
import { getApiUrl } from '../get-api-url';

export async function updateReview(review: UpdatedReview) {
  const response = await fetch(`/api/reviews`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(review),
  });

  if (!response.ok) {
    throw new Error('Falied to update review');
  }

  return response.json();
}

export async function deleteReview(id: number) {
  const response = await fetch(`/api/reviews/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error('Falied to delete review');
  }

  return response.json();
}

export async function addReview(review: NewReview) {
  const response = await fetch(`/api/reviews`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(review),
  });

  if (!response.ok) {
    throw new Error('Falied to create review');
  }

  return response.json();
}

export async function getReviewByRevieweeReviewer(
  revieweeId: number,
  reviewerId: number,
): Promise<Review> {
  console.log(
    `fetching review with revieweeId: ${revieweeId} and reviewerId: ${reviewerId}`,
  );

  const url = getApiUrl(`/api/reviews/reviewee/${revieweeId}/reviewer/${reviewerId}`);
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}
