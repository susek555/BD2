import { HistoryOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOffersTable from '@/app/ui/(offers-table)/generic-offer-table';
import SingleHistoryOffer from './history-offer-card';

export default async function OffersHistory({
  offers, // TODO: remove this prop
}: {
  offers: HistoryOffer[];
}) {
  return (
    <GenericOffersTable<HistoryOffer>
      offers={offers}
      ItemComponent={SingleHistoryOffer}
    />
  );
}
