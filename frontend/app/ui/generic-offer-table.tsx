import { BaseOffer, SearchParams } from '@/app/lib/definitions';
import Link from 'next/link';

interface GenericTableProps<T extends BaseOffer> {
  fetchFunction: (params: SearchParams) => Promise<T[]>;
  params: SearchParams;
  ItemComponent: React.ComponentType<{ offer: T }>;
  className?: string;
}

export default async function GenericOffersTable<T extends BaseOffer>({
  fetchFunction,
  params,
  ItemComponent,
  className = '',
}: GenericTableProps<T>) {
  const offers = await fetchFunction(params);

  return (
    <div className={`flex flex-col gap-4 ${className}`}>
      {offers.map((offer) => (
        <Link
          key={offer.id}
          href={`/offer/${offer.id}`}
          className={`block rounded-lg transition-shadow hover:shadow-lg ${
            offer.isAuction ? 'ring-2 ring-blue-500' : 'ring-2 ring-gray-500'
          }`}
        >
          <ItemComponent offer={offer} />
        </Link>
      ))}
    </div>
  );
}
