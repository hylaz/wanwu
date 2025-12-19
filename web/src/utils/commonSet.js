import { i18n } from '@/lang';

export const CHAT = 'chatflow';
export const WORKFLOW = 'workflow';
export const RAG = 'rag';
export const AGENT = 'agent';
export const AppType = {
  [WORKFLOW]: i18n.t('appSpace.workflow'),
  [CHAT]: i18n.t('appSpace.chat'),
  [RAG]: i18n.t('appSpace.rag'),
  [AGENT]: i18n.t('appSpace.agent'),
};
export const SafetyType = {
  Political: i18n.t('common.safetyType.political'),
  Revile: i18n.t('common.safetyType.revile'),
  Pornography: i18n.t('common.safetyType.pornography'),
  ViolentTerror: i18n.t('common.safetyType.violentTerror'),
  Illegal: i18n.t('common.safetyType.illegal'),
  InformationSecurity: i18n.t('common.safetyType.informationSecurity'),
  Other: i18n.t('common.safetyType.other'),
};
