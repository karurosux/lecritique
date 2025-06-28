import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
  return {
    restaurantId: params.id,
    questionnaireId: params.questionnaireId
  };
};