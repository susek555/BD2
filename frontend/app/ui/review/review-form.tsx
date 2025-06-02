'use client';

import {
  NewReview,
  Review,
  UpdatedReview,
} from '@/app/lib/definitions/reviews';
import { ChangeEvent, useEffect, useState } from 'react';

interface ReviewFormProps {
  review?: Review | null;
  revieweeId?: number;
  onSubmitSuccess?: () => void;
}

export default function ReviewForm({
  review = null,
  revieweeId,
  onSubmitSuccess,
}: ReviewFormProps) {
  const [rating, setRating] = useState(0);
  const [hoverRating, setHoverRating] = useState(0);
  const [description, setDescription] = useState('');
  const [isSubmitted, setIsSubmitted] = useState(false);

  const maxCharacters = 200;
  const isEditMode = review !== null;

  useEffect(() => {
    if (review) {
      setRating(review.rating || 0);
      setDescription(review.description || '');
    }
  }, [review]);

  const handleStarClick = (starValue: number) => {
    setRating(starValue);
  };

  const handleDescriptionChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    if (value.length <= maxCharacters) {
      setDescription(value);
    }
  };

  const handleSubmit = () => {
    if (rating > 0) {
      if (isEditMode && review) {
        const updatedReview: UpdatedReview = {
          id: review.id,
          rating,
          description: description.trim(),
        };
      } else {
        if (!revieweeId) {
          console.error('revieweeId is required for creating new reviews');
          return;
        }

        const newReview: NewReview = {
          revieweeId,
          rating,
          description: description.trim(),
        };
      }
    }
  };

  const handleDelete = () => {
    if (!review) return;

    if (window.confirm('Are you sure you want to delete this review?')) {
      console.log('Review deleted, id:', review.id);
    }
  };

  if (isSubmitted) {
    return (
      <div className='mx-auto max-w-md rounded-lg bg-white p-6 shadow-md'>
        <div className='text-center'>
          <div className='mb-4 text-4xl text-green-500'>✓</div>
          <h3 className='mb-2 text-lg font-semibold text-gray-900'>
            {isEditMode
              ? 'Review updated successfully!'
              : 'Thank you for your review!'}
          </h3>
          <p className='text-gray-600'>
            {isEditMode
              ? 'Your changes have been saved.'
              : 'Your feedback has been submitted successfully.'}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className='mx-auto max-w-md rounded-lg bg-white p-6 shadow-md'>
      <h2 className='mb-6 text-center text-2xl font-bold text-gray-900'>
        {isEditMode ? 'Edit Review' : 'Leave a Review'}
      </h2>

      <div className='space-y-6'>
        <div>
          <label className='mb-3 block text-sm font-medium text-gray-700'>
            Rating
          </label>
          <div className='flex space-x-1'>
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                type='button'
                onClick={() => handleStarClick(star)}
                onMouseEnter={() => setHoverRating(star)}
                onMouseLeave={() => setHoverRating(0)}
                className='focus:ring-opacity-50 rounded text-3xl transition-colors duration-150 focus:ring-2 focus:ring-yellow-400 focus:outline-none'
              >
                <span
                  className={
                    star <= (hoverRating || rating)
                      ? 'text-yellow-400'
                      : 'text-gray-300'
                  }
                >
                  ★
                </span>
              </button>
            ))}
          </div>
          {rating > 0 && (
            <p className='mt-2 text-sm text-gray-600'>
              You rated: {rating} star{rating !== 1 ? 's' : ''}
            </p>
          )}
        </div>

        <div>
          <label className='mb-2 block text-sm font-medium text-gray-700'>
            Description
          </label>
          <textarea
            value={description}
            onChange={handleDescriptionChange}
            placeholder='Share your experience...'
            rows={4}
            className='w-full resize-none rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-500 focus:outline-none'
          />
          <div className='mt-2 flex items-center justify-between'>
            <p className='text-sm text-gray-500'>
              Share your thoughts about your experience
            </p>
            <p
              className={`text-sm ${
                description.length >= maxCharacters * 0.9
                  ? 'text-red-500'
                  : 'text-gray-500'
              }`}
            >
              {description.length}/{maxCharacters}
            </p>
          </div>
        </div>

        {isEditMode ? (
          <div className='flex space-x-3'>
            <button
              onClick={handleDelete}
              className='focus:ring-opacity-50 flex-1 rounded-md bg-red-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-red-700 focus:ring-2 focus:ring-red-500 focus:outline-none'
            >
              Delete
            </button>
            <button
              onClick={handleSubmit}
              disabled={rating === 0 || description.trim() === ''}
              className='focus:ring-opacity-50 flex-1 rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-400'
            >
              Save Changes
            </button>
          </div>
        ) : (
          <button
            onClick={handleSubmit}
            disabled={rating === 0 || description.trim() === '' || !revieweeId}
            className='focus:ring-opacity-50 w-full rounded-md bg-blue-600 px-4 py-2 text-white transition-colors duration-200 hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-400'
          >
            Submit Review
          </button>
        )}

        {!isEditMode && !revieweeId && (
          <p className='text-center text-sm text-red-500'>
            revieweeId is required to create a new review
          </p>
        )}
      </div>
    </div>
  );
}
