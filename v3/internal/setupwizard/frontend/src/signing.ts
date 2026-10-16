import type { DarwinSigningDefaults } from './types';

function extractTeamID(identity: string): string {
  return identity.match(/\(([A-Z0-9]+)\)$/)?.[1] || '';
}

export function selectSigningIdentity(config: DarwinSigningDefaults | undefined, identity: string): DarwinSigningDefaults {
  const derivedTeamID = extractTeamID(identity);
  const previousDerivedTeamID = extractTeamID(config?.identity || '');
  const manualTeamID = config?.teamID !== previousDerivedTeamID ? config?.teamID : '';
  return { ...config, identity, teamID: derivedTeamID || manualTeamID || '' };
}
