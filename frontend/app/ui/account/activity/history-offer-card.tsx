import { HistoryOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOfferCard from '@/app/ui/(offers-table)/generic-offer-card';
import ReviewButton from './review-button';

export default function SingleHistoryOffer({ offer }: { offer: HistoryOffer }) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-EN', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    });
  };

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between border-b border-gray-200 pb-3'>
        <div className='flex space-x-8'>
          <div className='flex flex-col'>
            <p className='pt-4 pl-4 text-sm text-gray-600'>Purchased on</p>
            <p className='pl-4 font-bold'>{formatDate(offer.dateEnd)}</p>
          </div>
          <div className='flex flex-col'>
            <p className='pt-4 text-sm text-gray-600'>From</p>
            <p className='font-bold'>{offer.sellerName}</p>
          </div>
        </div>

        <div className='right pr-4'>
          <ReviewButton
            sellerId={offer.sellerId}
            sellerRating={offer.sellerRating}
          />
        </div>
      </div>
      <GenericOfferCard<HistoryOffer> offer={offer} variant='history' />
    </div>
  );
}
