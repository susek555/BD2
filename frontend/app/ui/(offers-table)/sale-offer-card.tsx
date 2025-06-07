'use client';

import { updateFavoriteStatus } from '@/app/lib/api/listing/requests';
import { SaleOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOfferCard, {
  GenericOfferProps,
} from '@/app/ui/(offers-table)/generic-offer-card';
import { StarIcon as StarOutline } from '@heroicons/react/24/outline';
import { StarIcon as StarSolid } from '@heroicons/react/24/solid';
import { useSession } from 'next-auth/react';
import { useState } from 'react';
import { AuthRequiredButton } from '../signup/auth-required-button';

export default function SaleOfferCard({
  offer,
  variant = 'default',
  priceLabel,
  className = '',
}: GenericOfferProps<SaleOffer>) {
  const [isFavorite, setIsFavorite] = useState(offer.is_liked);
  const { status } = useSession();
  const isLoggedIn = status === 'authenticated';

  const handleFavoriteToggle = () => {
    try{
      updateFavoriteStatus(offer.id, offer.is_liked);
      setIsFavorite(!isFavorite);
      offer.is_liked = !isFavorite;
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (error) {
      alert('Failed to update favorite status. Please try again later.');
    }
  };

  const favoriteButton = (
    <AuthRequiredButton
      isLoggedIn={isLoggedIn}
      onClick={handleFavoriteToggle}
      className='mr-2 ml-auto rounded-full p-1.5 transition-colors hover:bg-gray-100'
      aria-label={
        offer.is_liked ? 'Remove from favorites' : 'Add to favorites'
      }
    >
      {isFavorite ? (
        <StarSolid className='h-6 w-6 text-yellow-500' />
      ) : (
        <StarOutline className='h-6 w-6 text-gray-500 hover:text-yellow-400' />
      )}
    </AuthRequiredButton>
  );

  return (
    <GenericOfferCard
      offer={offer}
      variant={variant}
      priceLabel={priceLabel}
      headerContent={favoriteButton}
      className={className}
    />
  );
}
