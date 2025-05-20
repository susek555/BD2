import { BaseOffer } from '@/app/lib/definitions/SaleOffer';
import GenericOffersTable from '@/app/ui/generic-offer-table';
import SingleListingsOffer from './listings-offer-card';

export default async function UsersListingsTable({
  offers
}: {
  offers: BaseOffer[];
}) {
  return (
    <GenericOffersTable<BaseOffer>
      offers={offers}
      ItemComponent={SingleListingsOffer}
    />
  );
}
