{{- define "ingressAnnotations" }}
  {{- $annotations := dict }}
  {{- $componentSpec := index .Values .component }}

  {{- if hasKey $componentSpec.ingress "annotations" }}
    {{- $_ := merge $annotations $componentSpec.ingress.annotations }}
  {{- end }}

  {{- $_ := merge $annotations .Values.ingress.annotations .Values.jxRequirements.ingress.annotations  }}

  {{- if not (hasKey $annotations "kubernetes.io/ingress.class") }}
    {{- $customIngressClass := "" }}
    {{- if $componentSpec.ingress.customIngressClass }}
      {{- $customIngressClass := $componentSpec.ingress.customIngressClass }}
      {{- $_ := set $annotations "kubernetes.io/ingress.class" $customIngressClass }}
    {{- else if hasKey .Values.ingress "customIngressClass" }}
      {{- if eq .component "docker-registry" }}
        {{- if hasKey .Values.ingress.customIngressClass "dockerRegistry" }}
          {{- $customIngressClass := index .Values.ingress.customIngressClass "dockerRegistry" }}
          {{- $_ := set $annotations "kubernetes.io/ingress.class" $customIngressClass }}
        {{- end }}
      {{- else if hasKey .Values.ingress.customIngressClass .component }}
        {{- $customIngressClass := index .Values.ingress.customIngressClass .component }}
        {{- $_ := set $annotations "kubernetes.io/ingress.class" $customIngressClass }}
      {{- end }}
    {{- end }}
    {{- if not (hasKey $annotations "kubernetes.io/ingress.class") }}
      {{- $_ := set $annotations "kubernetes.io/ingress.class" ($customIngressClass | default "nginx")  }}
    {{- end }}
  {{- end }}
  {{- if and (hasKey .Values.jxRequirements.ingress "serviceType") (.Values.jxRequirements.ingress.serviceType) (eq .Values.jxRequirements.ingress.serviceType "NodePort") (not (hasKey $annotations "jenkins.io/host")) }}
    {{- $_ := set $annotations "jenkins.io/host" .Values.jxRequirements.ingress.domain }}
  {{- end }}
  {{- if $annotations }}
{{ toYaml $annotations | indent 4 }}
  {{- end }}
{{- end }}
