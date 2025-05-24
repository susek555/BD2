'use client';
import { RatingPercentages } from '@/app/lib/definitions/reviews';
import { useRef, useState } from 'react';
import { AverageRatingCard } from '../review/average-rating-card';
import { RatingDistribution } from '../review/rating-distribution';

interface RatingDisplayProps {
  rating: number;
  distribution: RatingPercentages;
}

export default function RatingDisplay({
  rating,
  distribution,
}: RatingDisplayProps) {
  const [showDistribution, setShowDistribution] = useState(false);
  const [isClosing, setIsClosing] = useState(false);
  const ratingRef = useRef<HTMLDivElement>(null);

  const handleClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      setShowDistribution(false);
      setIsClosing(false);
    }, 300); // animation duration
  };

  return (
    <>
      {showDistribution && (
        <div
          className={`fixed inset-0 z-10 ${isClosing ? 'animate-fadeOut' : 'animate-fadeIn'}`}
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
            if (!isClosing) {
              handleClose();
            }
          }}
        ></div>
      )}

      <div ref={ratingRef} className='relative z-20'>
        {!showDistribution ? (
          <div
            onClick={(e) => {
              setShowDistribution(true);
            }}
            className='cursor-pointer transition-all hover:opacity-80'
          >
            <AverageRatingCard rating={rating} />
          </div>
        ) : (
          <div
            className={`${isClosing ? 'animate-fadeOut' : 'animate-fadeIn'}`}
          >
            <div
              className='absolute -top-2 -right-2 flex h-6 w-6 cursor-pointer items-center justify-center rounded-full bg-gray-200 text-gray-600 hover:bg-gray-300'
              onClick={(e) => {
                e.stopPropagation();
                handleClose();
              }}
            >
              ×
            </div>
            <div className='w-64'>
              <RatingDistribution
                averageRating={rating}
                distribution={distribution}
              />
            </div>
          </div>
        )}
      </div>
    </>
  );
}
