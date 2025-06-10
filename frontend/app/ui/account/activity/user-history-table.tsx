import { cachedSessionData } from '@/app/lib/data/account/data';
import { HistoryOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOffersTable from '@/app/ui/(offers-table)/generic-offer-table';
import SingleHistoryOffer from './history-offer-card';

export default async function OffersHistory({
  offers, // TODO: remove this prop
}: {
  offers: HistoryOffer[];
}) {
  const userProfile = await cachedSessionData();

  const HistoryOfferWithUserId = (props: any) => (
    <SingleHistoryOffer {...props} userId={userProfile.userId} />
  );

  return (
    <GenericOffersTable<HistoryOffer>
      offers={offers}
      ItemComponent={HistoryOfferWithUserId}
    />
  );
}
