import { authConfig } from '@/app/lib/authConfig';
import { HistoryOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOffersTable from '@/app/ui/(offers-table)/generic-offer-table';
import { getServerSession } from 'next-auth';
import SingleHistoryOffer from './history-offer-card';

export default async function OffersHistory({
  offers, // TODO: remove this prop
}: {
  offers: HistoryOffer[];
}) {
  const session = await getServerSession(authConfig);

  const HistoryOfferWithUserId = (props: any) => (
    <SingleHistoryOffer {...props} userId={session?.user?.userId} />
  );

  return (
    <GenericOffersTable<HistoryOffer>
      offers={offers}
      ItemComponent={HistoryOfferWithUserId}
    />
  );
}
