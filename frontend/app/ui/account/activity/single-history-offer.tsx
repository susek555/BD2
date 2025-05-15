import { HistoryOffer } from '@/app/lib/definitions';
import GenericOffer from '@/app/ui/generic-offer';

export default function SingleHistoryOffer({ offer }: { offer: HistoryOffer }) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-EN', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    });
  };

  const headerContent = (
    <div className='pr-4 text-right'>
      <p className='text-sm text-gray-600'>Purchased on</p>
      <p className='font-bold'>{formatDate(offer.dateEnd)}</p>
    </div>
  );

  return (
    <GenericOffer<HistoryOffer>
      offer={offer}
      variant='history'
      headerContent={headerContent}
    />
  );
}
