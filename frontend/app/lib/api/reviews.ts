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
    throw new Error('Failed to update review');
  }

  return response;
}

export async function deleteReview(id: number) {
  const url = getApiUrl(`/api/reviews/${id}`);
  const response = await fetch(url, {
    method: 'DELETE',
  });

  console.log(response);

  if (!response.ok) {
    throw new Error('Failed to delete review');
  }

  return response;
}

export async function addReview(review: NewReview) {
  const response = await fetch(`/api/reviews`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      reviewee_id: review.revieweeId,
      description: review.description,
      rating: review.rating,
    }),
  });

  if (!response.ok) {
    throw new Error('Failed to create review');
  }

  return response;
}

export async function getReviewByRevieweeReviewer(
  revieweeId: number,
  reviewerId: number,
): Promise<Review | null> {
  console.log(
    `fetching review with revieweeId: ${revieweeId} and reviewerId: ${reviewerId}`,
  );

  const url = getApiUrl(
    `/api/reviews/reviewee/${revieweeId}/reviewer/${reviewerId}`,
  );
  const response = await fetch(url);

  if (!response.ok) {
    console.log(response);

    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const responseJson = await response.json();

  console.log(responseJson);

  if ('error_description' in responseJson) return null;

  return responseJson;
}
