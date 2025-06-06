import { NewReview, Review, UpdatedReview } from '../definitions/reviews';

export async function updateReview(review: UpdatedReview): Promise<Response> {
  console.log('Update review');

  return fetch(`/api/reviews`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(review),
  });
}

export async function deleteReview(id: number) {
  console.log('Delete review');

  return fetch(`/api/reviews/${id}`, {
    method: 'DELETE',
  });
}

export async function addReview(review: NewReview) {
  console.log('Add review');

  return fetch(`/api/reviews`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(review),
  });
}

export async function getReviewByRevieweeReviewer(
  revieweeId: number,
  reviewerId: number,
): Promise<Review> {
  console.log(
    `fetching review with revieweeId: ${revieweeId} and reviewerId: ${reviewerId}`,
  );

  const response = await fetch(
    `/api/reviews/reviewee/${revieweeId}/reviewer/${reviewerId}`,
  );

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}
