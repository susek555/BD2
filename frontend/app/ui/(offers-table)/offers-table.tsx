import { SaleOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOffersTable from './generic-offer-table';
import SaleOfferCard from './sale-offer-card';

export default async function OffersTable({
  offers,
}: {
  offers: SaleOffer[];
}) {
  return (
    <GenericOffersTable<SaleOffer>
      offers={offers}
      ItemComponent={SaleOfferCard}
    />
  );
}
