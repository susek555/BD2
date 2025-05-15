import { fetchHistory } from '@/app/lib/account/data';
import { SearchParams } from '@/app/lib/definitions';
import Link from 'next/link';
import SingleListingsOffer from './single-offer';

export default async function UsersListings({
  params,
}: {
  params: SearchParams;
}) {
  const offers = await fetchHistory(params);

  return (
    <div className='flex flex-col gap-4'>
      {offers.map((offer) => (
        <Link
          key={offer.id}
          href={`/account/listings/${offer.id}`} // TODO handle redirect to offer from users listings
          className={`block rounded-lg transition-shadow hover:shadow-lg ${
            offer.isAuction ? 'ring-2 ring-blue-500' : 'ring-2 ring-gray-500'
          }`}
        >
          <SingleListingsOffer offer={offer} />
        </Link>
      ))}
    </div>
  );
}
