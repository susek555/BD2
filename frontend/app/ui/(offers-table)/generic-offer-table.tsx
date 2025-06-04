import { BaseOffer } from '@/app/lib/definitions/SaleOffer';
import Link from 'next/link';

interface GenericTableProps<T extends BaseOffer> {
  offers: T[];
  ItemComponent: React.ComponentType<{ offer: T }>;
  className?: string;
}

export default async function GenericOffersTable<T extends BaseOffer>({
  offers,
  ItemComponent,
  className = '',
}: GenericTableProps<T>) {

  return (
    <div className={`flex flex-col gap-4 ${className}`}>
      {offers && offers.length > 0 ? (
      offers.map((offer) => (
        <Link
        key={offer.id}
        href={`/offer/${offer.id}`}
        className={`block rounded-lg transition-shadow hover:shadow-lg ${
          offer.is_auction ? 'ring-2 ring-blue-500' : 'ring-2 ring-gray-500'
        }`}
        >
        <ItemComponent offer={offer} />
        </Link>
      ))
      ) : (
      <div className="text-center py-6 text-gray-500 pr-20">
        No offers were found.
      </div>
      )}
    </div>
  );
}
