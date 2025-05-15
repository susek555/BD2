import { fetchHistory } from '@/app/lib/account/data';
import { HistoryOffer, SearchParams } from '@/app/lib/definitions';
import GenericOffersTable from '../../generic-offer-table';
import SingleHistoryOffer from './single-offer';

export default async function OffersHistory({
  params,
}: {
  params: SearchParams;
}) {
  return (
    <GenericOffersTable<HistoryOffer>
      fetchFunction={fetchHistory}
      params={params}
      ItemComponent={SingleHistoryOffer}
    />
  );
}
