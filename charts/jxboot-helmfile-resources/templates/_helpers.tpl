{{- define "ingressAnnotations" }}
  {{- $annotations := dict }}
  {{- $componentSpec := index .Values .component }}

  {{- if hasKey $componentSpec.ingress "annotations" }}
    {{- $_ := merge $annotations $componentSpec.ingress.annotations }}
  {{- end }}

  {{- $_ := merge $annotations .Values.ingress.annotations .Values.jxRequirements.ingress.annotations  }}

  {{- if and (hasKey .Values.jxRequirements.ingress "serviceType") (.Values.jxRequirements.ingress.serviceType) (eq .Values.jxRequirements.ingress.serviceType "NodePort") (not (hasKey $annotations "jenkins.io/host")) }}
    {{- $_ := set $annotations "jenkins.io/host" .Values.jxRequirements.ingress.domain }}
  {{- end }}
  {{- if $annotations }}
{{ toYaml $annotations | indent 4 }}
  {{- end }}
{{- end }}
